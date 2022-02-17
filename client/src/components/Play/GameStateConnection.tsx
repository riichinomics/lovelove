import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { ApiContext } from "../../rpc/ApiContext";
import { IState } from "../../state/IState";
import { Table } from "./Table";
import { useLocation, useNavigate } from "react-router";
import { CardMove, CardZone } from "../../utils";
import { CardMoveContext } from "../../rpc/CardMoveContext";
import { lovelove } from "../../rpc/proto/lovelove";
import { GameCreatedOnServerAction } from "../../state/actions/ConnectedToGameAction";
import { ActionType } from "../../state/actions/ActionType";
import { GameUpdateReceivedAction } from "../../state/actions/GameUpdateReceivedAction";
import { RoundEndClearedAction } from "../../state/actions/RoundEndClearedAction";
import { WaitingCurtain } from "./WaitingCurtain";
import { EndGameCurtain } from "./EndGameCurtain";
import { EnteredNewRoomAction } from "../../state/actions/EnteredNewRoomAction";
import { ApiConnection } from "../../rpc/Api";

export const GameStateConnection = () => {
	console.log("redraw");
	const { api } = React.useContext(ApiContext);

	const dispatch = useDispatch<React.Dispatch<GameCreatedOnServerAction | GameUpdateReceivedAction | RoundEndClearedAction | EnteredNewRoomAction>>();

	const roomId = useLocation().hash?.slice(1);
	const navigate = useNavigate();
	const userId = useSelector<IState, string>((state) => state.userId);
	const [apiConnection, setApiConnection] = React.useState<ApiConnection>(null);

	const position = useSelector((state: IState) => state.gamePosition);
	const opponentDisconnected = useSelector((state: IState) => state.opponentDisconnected);
	const gameState = useSelector((state: IState) => state.gameState);
	const roundEndView = useSelector((state: IState) => state.roundEndView);
	const [move, setMove] = React.useState<CardMove>();
	const [teyakuResolved, setTeyakuResolved] = React.useState(false);

	const [roomFull, setRoomFull] = React.useState(false);

	React.useEffect(() => {
		api.connect().then(setApiConnection);
	}, [api]);

	React.useEffect(() => {
		if (roomId == null || roomId === "") {
			navigate(
				{
					hash: Math.random().toString(16).substr(2, 5).toUpperCase()
				},
				{
					replace: true,
				}
			);
		}
	}, [roomId]);

	React.useEffect(() => {
		if (!apiConnection) {
			return;
		}

		apiConnection.closed.then(({reconnect}) => {
			setApiConnection(null);
			reconnect().then(setApiConnection);
		});

		apiConnection.lovelove.authenticate({
			userId,
		}).then((response) => {
			console.log("authentication response", response);
		});

		console.log(apiConnection);

		const messageSub = apiConnection.broadcastMessages.subscribe(message => {
			console.log(message);

			if (message.constructor.name !== "GameStateUpdate") {
				return;
			}

			setMove(null);
			setTeyakuResolved(false);

			const gameStateUpdate = message as any as lovelove.GameStateUpdate;
			dispatch({
				type: ActionType.GameUpdateReceived,
				update: gameStateUpdate
			});
		});

		return () => {
			messageSub.unsubscribe();
		};
	}, [apiConnection]);

	React.useEffect(() => {
		if (!apiConnection) {
			return;
		}

		console.log("requesting GameState");
		setRoomFull(false);
		dispatch({
			type: ActionType.EnteredNewRoom
		});

		apiConnection.lovelove.connectToGame({
			roomId
		}).then(response => {
			console.log("GameStateConnection", response);
			if (response.status === lovelove.ConnectToGameResponseCode.ConnectToGameFull) {
				dispatch({
					type: ActionType.ConnectedToGame,
					position: lovelove.PlayerPosition.UnknownPosition,
					opponentDisconnected: false,
				});
				setRoomFull(true);
				return;
			}

			dispatch({
				type: ActionType.ConnectedToGame,
				position: response.playerPosition,
				opponentDisconnected: response.OpponentDisconnected,
			});
		});
	}, [dispatch, apiConnection, roomId]);

	React.useEffect(() => {
		if (!apiConnection) {
			return;
		}

		if (!move) {
			return;
		}

		if (move.from.zone === CardZone.Hand) {
			const request: lovelove.IPlayHandCardRequest = {
				handCard: {
					cardId: move.from.card.id
				},
			};

			if (move.to.card) {
				request.tableCard = {
					cardId: move.to.card.id
				};
			}

			apiConnection.lovelove.playHandCard(request).then(response => {
				if (response.status === lovelove.GenericResponseCode.Error) {
					setMove(null);
				}
			});

			return;
		}

		if (move.from.zone === CardZone.Drawn) {
			const request: lovelove.IPlayDrawnCardRequest = {
				tableCard: {
					cardId: move.to.card.id
				}
			};

			apiConnection.lovelove.playDrawnCard(request).then(response => {
				if (response.status === lovelove.GenericResponseCode.Error) {
					setMove(null);
				}
			});

			return;
		}


	}, [move, setMove, apiConnection]);

	const onCardDropped = React.useCallback((move: CardMove) => {
		console.log(move);
		setMove(move);
	}, [setMove]);


	const onKoikoiChosen = React.useCallback(() => {
		apiConnection?.lovelove?.resolveShoubuOpportunity({});
	}, [apiConnection]);

	const onShoubuChosen = React.useCallback((teyaku: boolean) => {
		if (!apiConnection) {
			return;
		}

		if (teyaku) {
			apiConnection.lovelove.resolveTeyaku({}).then((response) => {
				if (response.status != lovelove.GenericResponseCode.Error) {
					setTeyakuResolved(true);
				}
			});
			return;
		}
		apiConnection.lovelove.resolveShoubuOpportunity({shoubu: true});
	}, [apiConnection]);

	const onContinueChosen = React.useCallback(() => {
		dispatch({
			type: ActionType.RoundEndCleared,
		});
	}, [dispatch]);


	const onGameConceded = React.useCallback(() => {
		apiConnection?.lovelove?.concedeGame({});
	}, [apiConnection]);

	const onRematchRequested = React.useCallback(() => {
		apiConnection?.lovelove?.requestRematch({});
	}, [apiConnection]);

	if (!gameState || !apiConnection) {
		return <WaitingCurtain roomFull={roomFull} connected={!!apiConnection} />;
	}

	if (gameState?.gameEnd) {
		return <EndGameCurtain
			position={position}
			gameState={gameState}
			onRematchRequested={onRematchRequested}
		/>;
	}

	return <CardMoveContext.Provider value={{move}}>
		<Table
			opponentDisconnected={opponentDisconnected}
			gameState={roundEndView?.gameState ?? gameState}
			position={position}
			onCardDropped={onCardDropped}
			onKoikoiChosen={onKoikoiChosen}
			onShoubuChosen={onShoubuChosen}
			roundEndView={roundEndView}
			onContinueChosen={onContinueChosen}
			teyakuResolved={teyakuResolved}
			onGameConceded={onGameConceded}
		/>
	</CardMoveContext.Provider>;
};
