import * as React from "react";
import * as ReactDOM from "react-dom";
import { compose, createStore } from "redux";
import { ActionType } from "./state/actions/ActionType";
import { Api } from "./rpc/Api";
import { ApiContext } from "./rpc/ApiContext";
import { ApiStateChangedAction } from "./state/actions/ApiStateChangedAction";
import { GameStateConnection } from "./components/Play/GameStateConnection";
import { IState } from "./state/IState";
import { Provider } from "react-redux";
import { ThemeContext } from "./themes/ThemeContext";
import mainReducer from "./state/reducers";
import mantia from "./themes/mantia";
import { ApiState } from "./rpc/ApiState";
import * as uuid from "uuid";

const USER_ID_LOCAL_STORAGE_KEY = "lovelove_user_id";
const storedUserId = localStorage.getItem(USER_ID_LOCAL_STORAGE_KEY);
const userId = storedUserId ?? uuid.v4();

if (!storedUserId) {
	localStorage.setItem(USER_ID_LOCAL_STORAGE_KEY, userId);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const composeEnhancers = (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = createStore(
	mainReducer,
	{
		userId,
		apiState: ApiState.Connecting
	} as IState,
	composeEnhancers()
);

const api = new Api({url: "localhost:6482"});
api.init().then(() => {
	api.lovelove.authenticate({
		userId,
	}).then((response) => {
		console.log("authentication response", response);
		store.dispatch<ApiStateChangedAction>({
			type: ActionType.ApiStateChanged,
			apiState: ApiState.Connected
		});
	});
});

ReactDOM.render(
	<Provider store={store}>
		<ApiContext.Provider value={{api}}>
			<ThemeContext.Provider value={{
				theme: mantia
			}}
			>
				{/* <Table
					collection={[...Array(8 * 4)].map(_ => createRandomCard())}
					// drawnCard={createRandomCard()}
					deck={Math.random() * 4 | 0}
					hand={[...Array(Math.random() * 8 | 0)].map(_ => createRandomCard())}
					opponentCollection={[...Array(8 * 4)].map(_ => createRandomCard())}
					opponentHand={Math.random() * 8 | 0}
					table={[...Array(12 + Math.random() * 6 | 0)].map(_ => createRandomCard())}
				/> */}
				<GameStateConnection />
			</ThemeContext.Provider>
		</ApiContext.Provider>
	</Provider>,
	document.getElementsByTagName("body")[0]
);
