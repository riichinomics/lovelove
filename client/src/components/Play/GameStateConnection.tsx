import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { ActionType } from "../../state/actions/ActionType";
import { ApiContext } from "../../rpc/ApiContext";
import { IState } from "../../state/IState";
import { InitialGameStateReceivedAction } from "../../state/actions/InitialGameStateReceivedAction";
import { Table } from "./Table";
import { ApiState } from "../../rpc/ApiState";
import { useLocation, useNavigate } from "react-router";
import { CardMove, createRandomCard } from "./utils";
import { CardMoveContext } from "../../rpc/CardMoveContext";

export const GameStateConnection = () => {
	const { api } = React.useContext(ApiContext);

	const dispatch = useDispatch<React.Dispatch<InitialGameStateReceivedAction>>();

	const roomId = useLocation().hash?.slice(1);
	const navigate = useNavigate();
	const apiState = useSelector<IState>((state) => state.apiState);

	React.useEffect(() => {
		if (roomId == null || roomId === "") {
			navigate({
				hash: Math.random().toString(16).substr(2, 5).toUpperCase()
			});
		}
	}, [roomId]);

	React.useEffect(() => {
		if (roomId === "test") {
			dispatch({
				type: ActionType.InitialGameStateReceived,
				position: Math.random() * 2 | 0,
				gameState: {
					collection: [...Array(8 * 4)].map(_ => createRandomCard()),
					//: drawnCard={createRandomCard()},
					deck: Math.random() * 4 | 0,
					hand: [...Array(Math.random() * 8 | 0)].map(_ => createRandomCard()),
					opponentCollection: [...Array(8 * 4)].map(_ => createRandomCard()),
					opponentHand: Math.random() * 8 | 0,
					table: [...Array(12 + Math.random() * 6 | 0)].map(_ => createRandomCard()),
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
	}, [api, dispatch, apiState, roomId]);

	const onCardDropped = React.useCallback((move: CardMove) => {
		console.log(move);
		setMove(move);
	}, []);

	const [move, setMove] = React.useState<CardMove>();

	const gameState = useSelector((state: IState) => state.gameState ?? {});
	const position = useSelector((state: IState) => state.gamePosition);
	return <CardMoveContext.Provider value={{move}}>
		<Table
			{...gameState}
			position={position}
			onCardDropped={onCardDropped}
		/>
	</CardMoveContext.Provider>;
};
