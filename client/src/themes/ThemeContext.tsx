import * as React from "react";
import { ITheme } from "./ITheme";
import mantia from "../themes/mantia";

interface ThemeContextProps {
	theme: ITheme;
}

export const ThemeContext = React.createContext<ThemeContextProps>({
	theme: mantia
});
