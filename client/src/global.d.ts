declare module "*.sass" {
	interface IClassNames {
		[className: string]: string
	}
	const classNames: IClassNames;
	export = classNames;
}

declare module "*.scss" {
	interface IClassNames {
		[className: string]: string
	}
	const classNames: IClassNames;
	export = classNames;
}

declare module "*.svg" {
	const ReactComponent: React.FC<React.SVGProps<SVGSVGElement>>;
	const content: string;

	export { ReactComponent };
	export default content;
 }
