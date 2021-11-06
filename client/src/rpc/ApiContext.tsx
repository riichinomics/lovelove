import * as React from "react";
import { Api } from "./Api";

interface ApiContextProps {
	api?: Api;
}

export const ApiContext = React.createContext<ApiContextProps>({});
