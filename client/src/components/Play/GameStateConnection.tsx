import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { ActionType } from "../../state/actions/ActionType";
import { ApiContext } from "../../rpc/ApiContext";
import { IState } from "../../state/IState";
import { InitialGameStateReceivedAction } from "../../state/actions/InitialGameStateReceivedAction";
import { Table } from "./Table";
import { ApiState } from "../../rpc/ApiState";

export const GameStateConnection = () => {
	const { api } = React.useContext(ApiContext);
	const dispatch = useDispatch<React.Dispatch<InitialGameStateReceivedAction>>();

	const apiState = useSelector<IState>((state) => state.apiState);

	React.useEffect(() => {
		if (apiState !== ApiState.Connected) {
			return;
		}

		console.log("requesting GameState");

		api.lovelove.connectToGame({
			roomId: "roomId"
		}).then(response => {
			console.log("GameStateConnection", response);
			dispatch({
				type: ActionType.InitialGameStateReceived,
				gameState: response.gameState
			});
		});
	}, [api, dispatch, apiState]);

	const gameState = useSelector<IState>((state) => state.gameState ?? {});
	return <Table {...gameState} />;
};
