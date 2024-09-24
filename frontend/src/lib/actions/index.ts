import type { Action } from 'svelte/action';

export const selectTextOnFocus: Action<HTMLInputElement> = (node) => {
	const handleFocus = () => {
		if (node && typeof node.select === 'function') {
			node.select();
		}
	};

	node.addEventListener('focus', handleFocus);

	return {
		destroy() {
			node.removeEventListener('focus', handleFocus);
		}
	};
};
