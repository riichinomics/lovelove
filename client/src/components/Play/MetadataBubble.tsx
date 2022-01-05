import { stylesheet } from "astroturf";
import clsx from "clsx";
import * as React from "react";

const styles = stylesheet`
	.metadataBubble {
		border-radius: 3px;
		background-color: #333;
		min-width: 49px;
		box-sizing: border-box;
		text-align: center;
		font-weight: 800;
		padding: 6px 12px 8px;
		font-size: 24px;
	}
`;

export const MetadataBubble: React.FC<{
	className?: string;
}> = (props) => {
	return <div className={clsx(styles.metadataBubble, props.className)}>
		{props.children}
	</div>;
};
