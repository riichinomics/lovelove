/* eslint-disable sort-imports */
import * as React from "react";
import { CardProps } from "../CardProps";
import { stylesheet } from "astroturf";

import ayame1 from "../mantia/ayame/1.svg";
import ayame2 from "../mantia/ayame/2.svg";
import ayame3 from "../mantia/ayame/3.svg";
import ayame4 from "../mantia/ayame/4.svg";

import botan1 from "../mantia/botan/1.svg";
import botan2 from "../mantia/botan/2.svg";
import botan3 from "../mantia/botan/3.svg";
import botan4 from "../mantia/botan/4.svg";

import fuji1 from "../mantia/fuji/1.svg";
import fuji2 from "../mantia/fuji/2.svg";
import fuji3 from "../mantia/fuji/3.svg";
import fuji4 from "../mantia/fuji/4.svg";

import hagi1 from "../mantia/hagi/1.svg";
import hagi2 from "../mantia/hagi/2.svg";
import hagi3 from "../mantia/hagi/3.svg";
import hagi4 from "../mantia/hagi/4.svg";

import kiku1 from "../mantia/kiku/1.svg";
import kiku2 from "../mantia/kiku/2.svg";
import kiku3 from "../mantia/kiku/3.svg";
import kiku4 from "../mantia/kiku/4.svg";

import kiri1 from "../mantia/kiri/1.svg";
import kiri2 from "../mantia/kiri/2.svg";
import kiri3 from "../mantia/kiri/3.svg";
import kiri4 from "../mantia/kiri/4.svg";

import matsu1 from "../mantia/matsu/1.svg";
import matsu2 from "../mantia/matsu/2.svg";
import matsu3 from "../mantia/matsu/3.svg";
import matsu4 from "../mantia/matsu/4.svg";

import momiji1 from "../mantia/momiji/1.svg";
import momiji2 from "../mantia/momiji/2.svg";
import momiji3 from "../mantia/momiji/3.svg";
import momiji4 from "../mantia/momiji/4.svg";

import sakura1 from "../mantia/sakura/1.svg";
import sakura2 from "../mantia/sakura/2.svg";
import sakura3 from "../mantia/sakura/3.svg";
import sakura4 from "../mantia/sakura/4.svg";

import susuki1 from "../mantia/susuki/1.svg";
import susuki2 from "../mantia/susuki/2.svg";
import susuki3 from "../mantia/susuki/3.svg";
import susuki4 from "../mantia/susuki/4.svg";

import ume1 from "../mantia/ume/1.svg";
import ume2 from "../mantia/ume/2.svg";
import ume3 from "../mantia/ume/3.svg";
import ume4 from "../mantia/ume/4.svg";

import yanagi1 from "../mantia/yanagi/1.svg";
import yanagi2 from "../mantia/yanagi/2.svg";
import yanagi3 from "../mantia/yanagi/3.svg";
import yanagi4 from "../mantia/yanagi/4.svg";

import { lovelove } from "../../rpc/proto/lovelove";

const SVG_MAP: Record<lovelove.Hana, Record<lovelove.Variation, string>> = {
	[lovelove.Hana.UnknownSeason]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: "",
		[lovelove.Variation.Second]: "",
		[lovelove.Variation.Third]: "",
		[lovelove.Variation.Fourth]: "",
	},
	[lovelove.Hana.Ayame]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: ayame1,
		[lovelove.Variation.Second]: ayame2,
		[lovelove.Variation.Third]: ayame3,
		[lovelove.Variation.Fourth]: ayame4,
	},
	[lovelove.Hana.Botan]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: botan1,
		[lovelove.Variation.Second]: botan2,
		[lovelove.Variation.Third]: botan3,
		[lovelove.Variation.Fourth]: botan4,
	},
	[lovelove.Hana.Fuji]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: fuji1,
		[lovelove.Variation.Second]: fuji2,
		[lovelove.Variation.Third]: fuji3,
		[lovelove.Variation.Fourth]: fuji4,
	},
	[lovelove.Hana.Hagi]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: hagi1,
		[lovelove.Variation.Second]: hagi2,
		[lovelove.Variation.Third]: hagi3,
		[lovelove.Variation.Fourth]: hagi4,
	},
	[lovelove.Hana.Kiku]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: kiku1,
		[lovelove.Variation.Second]: kiku2,
		[lovelove.Variation.Third]: kiku3,
		[lovelove.Variation.Fourth]: kiku4,
	},
	[lovelove.Hana.Kiri]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: kiri1,
		[lovelove.Variation.Second]: kiri2,
		[lovelove.Variation.Third]: kiri3,
		[lovelove.Variation.Fourth]: kiri4,
	},
	[lovelove.Hana.Matsu]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: matsu1,
		[lovelove.Variation.Second]: matsu2,
		[lovelove.Variation.Third]: matsu3,
		[lovelove.Variation.Fourth]: matsu4,
	},
	[lovelove.Hana.Momiji]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: momiji1,
		[lovelove.Variation.Second]: momiji2,
		[lovelove.Variation.Third]: momiji3,
		[lovelove.Variation.Fourth]: momiji4,
	},
	[lovelove.Hana.Sakura]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: sakura1,
		[lovelove.Variation.Second]: sakura2,
		[lovelove.Variation.Third]: sakura3,
		[lovelove.Variation.Fourth]: sakura4,
	},
	[lovelove.Hana.Susuki]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: susuki1,
		[lovelove.Variation.Second]: susuki2,
		[lovelove.Variation.Third]: susuki3,
		[lovelove.Variation.Fourth]: susuki4,
	},
	[lovelove.Hana.Ume]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: ume1,
		[lovelove.Variation.Second]: ume2,
		[lovelove.Variation.Third]: ume3,
		[lovelove.Variation.Fourth]: ume4,
	},
	[lovelove.Hana.Yanagi]: {
		[lovelove.Variation.UnknownVariation]: "",
		[lovelove.Variation.First]: yanagi1,
		[lovelove.Variation.Second]: yanagi2,
		[lovelove.Variation.Third]: yanagi3,
		[lovelove.Variation.Fourth]: yanagi4,
	}
};

const styles = stylesheet`
	.card {
		width: 100px;
		> img {
			display: block;

		}
	}
`;

export const CardComponent: React.FC<CardProps> = (props) => {
	return <div className={styles.card}>
		<img draggable={false} src={SVG_MAP[props.hana][props.variation]} />
	</div>;
};
