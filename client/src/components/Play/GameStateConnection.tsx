import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { ApiContext } from "../../rpc/ApiContext";
import { IState } from "../../state/IState";
import { Table } from "./Table";
import { ApiState } from "../../rpc/ApiState";
import { useLocation, useNavigate } from "react-router";
import { CardMove, CardZone, createRandomCard } from "./utils";
import { CardMoveContext } from "../../rpc/CardMoveContext";
import { lovelove } from "../../rpc/proto/lovelove";
import { InitialGameStateReceivedAction } from "../../state/actions/InitialGameStateReceivedAction";
import { ActionType } from "../../state/actions/ActionType";
import { GameUpdateReceivedAction } from "../../state/actions/GameUpdateReceivedAction";


export const GameStateConnection = () => {
	const { api } = React.useContext(ApiContext);

	const dispatch = useDispatch<React.Dispatch<InitialGameStateReceivedAction | GameUpdateReceivedAction>>();

	const roomId = useLocation().hash?.slice(1);
	const navigate = useNavigate();
	const apiState = useSelector<IState>((state) => state.apiState);

	const position = useSelector((state: IState) => state.gamePosition);
	const gameState = useSelector((state: IState) => state.gameState);
	const [move, setMove] = React.useState<CardMove>();

	React.useEffect(() => {
		if (roomId == null || roomId === "") {
			navigate({
				hash: Math.random().toString(16).substr(2, 5).toUpperCase()
			});
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
		if (roomId === "test") {
			dispatch({
				type: ActionType.InitialGameStateReceived,
				position: Math.random() * 2 | 0,
				gameState: {
					collection: [...Array(8 * 4)].map(() => createRandomCard()),
					//: drawnCard={createRandomCard()},
					deck: Math.random() * 4 | 0,
					hand: [...Array(Math.random() * 8 | 0)].map(() => createRandomCard()),
					opponentCollection: [...Array(8 * 4)].map(() => createRandomCard()),
					opponentHand: Math.random() * 8 | 0,
					table: [...Array(12 + Math.random() * 6 | 0)].map(() => ({
						card: createRandomCard()
					})),
					oya: Math.random() * 2 | 0,
					active: Math.random() * 2 | 0,
				}
			});
			return;
		}

		if (apiState !== ApiState.Connected) {
			return;
		}

		console.log("requesting GameState");

		api.lovelove.connectToGame({
			roomId
		}).then(response => {
			console.log("GameStateConnection", response);

			dispatch({
				type: ActionType.InitialGameStateReceived,
				position: response.position,
				gameState: response.gameState
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


	const koikoiChosen  = React.useCallback(() => {
		api.lovelove.resolveShoubuOpportunity({});
	}, [api]);

	const shoubuChosen  = React.useCallback(() => {
		api.lovelove.resolveShoubuOpportunity({shoubu: true});
	}, [api]);

	return <CardMoveContext.Provider value={{move}}>
		<Table
			{...gameState}
			position={position}
			onCardDropped={onCardDropped}
			koikoiChosen={koikoiChosen}
			shoubuChosen={shoubuChosen}
		/>
	</CardMoveContext.Provider>;
};
