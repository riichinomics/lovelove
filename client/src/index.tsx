import { Api } from "./rpc/Api";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { createStore, compose, Action } from "redux";
import { Provider } from "react-redux";
import { Play } from "./components/play";

enum ActionType {

}

interface IState {

}

const api = new Api({url: "localhost:6482"});
api.init().then(() => {
	api.lovelove.sayHello({
		name: "test"
	}).then((response) => {
		console.log("test", response)
	});
});

function mainReducer(state: IState, action: Action<ActionType>): IState {
	switch (action.type) {
	}

	return state;
}

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
		<Play.Table/>
	</Provider>,
	document.getElementsByTagName("body")[0]
);
