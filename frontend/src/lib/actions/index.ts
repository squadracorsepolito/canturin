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

export const focusOnDisplay: Action<HTMLElement> = (el) => {
	el.focus();
};

export function clickOutside(el: HTMLElement, handler: () => void) {
	const onClick = (e: MouseEvent) =>
		el && !el.contains(e.target as HTMLElement) && !e.defaultPrevented && handler();

	document.addEventListener('click', onClick, true);

	return {
		destroy() {
			document.removeEventListener('click', onClick, true);
		}
	};
}
