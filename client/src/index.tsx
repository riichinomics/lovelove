import * as React from "react";
import * as ReactDOM from "react-dom";
import { compose, createStore } from "redux";
import { Api } from "./rpc/Api";
import { ApiContext } from "./rpc/ApiContext";
import { GameStateConnection } from "./components/Play/GameStateConnection";
import { IState } from "./state/IState";
import { Provider } from "react-redux";
import { ThemeContext } from "./themes/ThemeContext";
import mainReducer from "./state/reducers";
import mantia from "./themes/mantia";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import * as uuid from "uuid";

import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";
import { InfoLinksWrapper } from "./components/InfoLinksWrapper";

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
	} as IState,
	composeEnhancers()
);

const api = new Api({
	production: location.protocol === "https:",
	url: window.location.hostname,
});
const ApiInitialiser: React.FC = ({children}) => {
	const {api} = React.useContext(ApiContext);
	const [initialised, setInitialised] = React.useState(false);
	React.useEffect(() => {
		setInitialised(false);
		api.init().then(() => {
			setInitialised(true);
		});
	}, [api]);

	if (!initialised) {
		return null;
	}

	// eslint-disable-next-line react/jsx-no-useless-fragment
	return <>{children}</>;
};


ReactDOM.render(
	<Provider store={store}>
		<ApiContext.Provider value={{api}}>
			<ApiInitialiser>
				<ThemeContext.Provider value={{
					theme: mantia
				}}
				>
					<DndProvider backend={HTML5Backend}>
						<BrowserRouter>
							<Routes>
								<Route path="/" element={<InfoLinksWrapper><GameStateConnection /></InfoLinksWrapper>} />
							</Routes>
						</BrowserRouter>
					</DndProvider>
				</ThemeContext.Provider>
			</ApiInitialiser>
		</ApiContext.Provider>
	</Provider>,
	document.getElementsByTagName("body")[0]
);
