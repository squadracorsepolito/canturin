const privateItemKey = Symbol('Item');

export type HighlightState = 'none' | 'selecting' | 'moving';

export type Item = {
	[privateItemKey]: true;
	instanceId: string;
	id: string;
};

export function getItem(data: Omit<Item, typeof privateItemKey>) {
	return {
		[privateItemKey]: true,
		...data
	};
}

export function isItem(data: Record<string | symbol, unknown>): data is Item {
	return Boolean(data[privateItemKey]);
}
