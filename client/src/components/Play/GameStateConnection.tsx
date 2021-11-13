import * as React from "react";
import { useDispatch, useSelector } from "react-redux";
import { ActionType } from "../../state/actions/ActionType";
import { ApiContext } from "../../rpc/ApiContext";
import { IState } from "../../state/IState";
import { InitialGameStateReceivedAction } from "../../state/actions/InitialGameStateReceivedAction";
import { Table } from "./Table";
import { ApiState } from "../../rpc/ApiState";
import { useLocation, useNavigate } from "react-router";

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

	const gameState = useSelector((state: IState) => state.gameState ?? {});
	const position = useSelector((state: IState) => state.gamePosition);
	return <Table {...gameState} position={position} />;
};
