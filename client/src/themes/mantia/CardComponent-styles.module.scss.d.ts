declare namespace CardComponentStylesModuleScssNamespace {
  export interface ICardComponentStylesModuleScss {
    card: string;
  }
}

declare const CardComponentStylesModuleScssModule: CardComponentStylesModuleScssNamespace.ICardComponentStylesModuleScss & {
  /** WARNING: Only available when `css-loader` is used without `style-loader` or `mini-css-extract-plugin` */
  locals: CardComponentStylesModuleScssNamespace.ICardComponentStylesModuleScss;
};

export = CardComponentStylesModuleScssModule;
