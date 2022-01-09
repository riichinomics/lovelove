import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { ApiContext } from "../../rpc/ApiContext";
import { IState } from "../../state/IState";
import { Table } from "./Table";
import { ApiState } from "../../rpc/ApiState";
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


export const GameStateConnection = () => {
	const { api } = React.useContext(ApiContext);

	const dispatch = useDispatch<React.Dispatch<GameCreatedOnServerAction | GameUpdateReceivedAction | RoundEndClearedAction | EnteredNewRoomAction>>();

	const roomId = useLocation().hash?.slice(1);
	const navigate = useNavigate();
	const apiState = useSelector<IState>((state) => state.apiState);

	const position = useSelector((state: IState) => state.gamePosition);
	const opponentDisconnected = useSelector((state: IState) => state.opponentDisconnected);
	const gameState = useSelector((state: IState) => state.gameState);
	const roundEndView = useSelector((state: IState) => state.roundEndView);
	const [move, setMove] = React.useState<CardMove>();
	const [teyakuResolved, setTeyakuResolved] = React.useState(false);

	const [roomFull, setRoomFull] = React.useState(false);

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
		if (apiState !== ApiState.Connected) {
			return;
		}

		const messageSub = api.broadcastMessages.subscribe(message => {
			console.log(message);

			switch (message.$type.name) {
				case lovelove.GameStateUpdate.name: {
					setMove(null);
					setTeyakuResolved(false);

					const gameStateUpdate = message as any as lovelove.GameStateUpdate;
					dispatch({
						type: ActionType.GameUpdateReceived,
						update: gameStateUpdate
					});
					break;
				}
			}
		});

		return () => {
			messageSub.unsubscribe();
		};
	}, [api, apiState]);

	React.useEffect(() => {
		if (apiState !== ApiState.Connected) {
			return;
		}

		console.log("requesting GameState");
		setRoomFull(false);
		dispatch({
			type: ActionType.EnteredNewRoom
		});

		api.lovelove.connectToGame({
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
	}, [dispatch, api, apiState, roomId]);

	React.useEffect(() => {
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

			api.lovelove.playHandCard(request).then(response => {
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

			api.lovelove.playDrawnCard(request).then(response => {
				if (response.status === lovelove.GenericResponseCode.Error) {
					setMove(null);
				}
			});

			return;
		}


	}, [move, setMove]);

	const onCardDropped = React.useCallback((move: CardMove) => {
		console.log(move);
		setMove(move);
	}, [setMove]);


	const onKoikoiChosen = React.useCallback(() => {
		api.lovelove.resolveShoubuOpportunity({});
	}, [api]);

	const onShoubuChosen = React.useCallback((teyaku: boolean) => {
		if (teyaku) {
			api.lovelove.resolveTeyaku({}).then((response) => {
				if (response.status != lovelove.GenericResponseCode.Error) {
					setTeyakuResolved(true);
				}
			});
			return;
		}
		api.lovelove.resolveShoubuOpportunity({shoubu: true});
	}, [api]);

	const onContinueChosen = React.useCallback(() => {
		dispatch({
			type: ActionType.RoundEndCleared,
		});
	}, [dispatch]);


	const onGameConceded = React.useCallback(() => {
		api.lovelove.concedeGame({});
	}, [api]);

	const onRematchRequested = React.useCallback(() => {
		api.lovelove.requestRematch({});
	}, [api]);

	if (gameState?.gameEnd) {
		return <EndGameCurtain
			position={position}
			gameState={gameState}
			onRematchRequested={onRematchRequested}
		/>;
	}

	if (!gameState) {
		return <WaitingCurtain roomFull={roomFull} />;
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
