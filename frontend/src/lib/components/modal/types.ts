import type { Snippet } from 'svelte';
import type { HTMLButtonAttributes } from 'svelte/elements';

export type ModalProps<T> = {
	onsubmit: (args: T) => void;
	trigger: Snippet<[{ getProps: () => HTMLButtonAttributes }]>;
};
