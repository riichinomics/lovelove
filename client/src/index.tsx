import * as React from "react";
import * as ReactDOM from "react-dom";
import { Action, compose, createStore } from "redux";
import { Api } from "./rpc/Api";
import { Provider } from "react-redux";
import { Table } from "./components/Play/Table";
import { ThemeContext } from "./themes/ThemeContext";
import mantia from "./themes/mantia";
import { createRandomCard } from "./components/Play/utils";

enum ActionType {

}

// eslint-disable-next-line @typescript-eslint/no-empty-interface
interface IState {

}

const api = new Api({url: "localhost:6482"});
api.init().then(() => {
	api.lovelove.sayHello({
		name: "test"
	}).then((response) => {
		console.log("test", response);
	});
});

function mainReducer(state: IState, action: Action<ActionType>): IState {
	// eslint-disable-next-line no-empty
	switch (action.type) {
	}

	return state;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const composeEnhancers = (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

const store = createStore(
	mainReducer,
	{
		contestsByMajsoulFriendlyId: {},
		musicPlayer: {
			playing: false,
			videoId: null
		},
	} as IState,
	composeEnhancers()
);

ReactDOM.render(
	<Provider store={store}>
		<ThemeContext.Provider value={{
			theme: mantia
		}}
		>
			<Table
				deck={Math.random() * 4 | 0}
				drawnCard={createRandomCard()}
				opponentCards={Math.random() * 8 | 0}
				opponentCollection={[...Array(8 * 4)].map(_ => createRandomCard())}
				playerCollection={[...Array(8 * 4)].map(_ => createRandomCard())}
				playerHand={[...Array(Math.random() * 8 | 0)].map(_ => createRandomCard())}
				sharedCards={[...Array(12 + Math.random() * 6 | 0)].map(_ => createRandomCard())}
			/>
		</ThemeContext.Provider>
	</Provider>,
	document.getElementsByTagName("body")[0]
);
