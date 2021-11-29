import * as React from "react";
import { CardMove } from "../components/Play/utils";

interface CardMoveContext {
	move?: CardMove;
}

export const CardMoveContext = React.createContext<CardMoveContext>({});
