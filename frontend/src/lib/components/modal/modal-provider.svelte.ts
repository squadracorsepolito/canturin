import * as dialog from '@zag-js/dialog';
import type { HTMLButtonAttributes } from 'svelte/elements';

const provider = new Map<string, () => dialog.Api>();

export function getModalTrigger(key: string) {
	const apiGetter = provider.get(key);

	if (!apiGetter) {
		throw new Error(`Modal with key ${key} not registered`);
	}

	return apiGetter().getTriggerProps() as HTMLButtonAttributes;
}

export function openModal(key: string) {
	const apiGetter = provider.get(key);

	if (!apiGetter) {
		throw new Error(`Modal with key ${key} not registered`);
	}

	apiGetter().setOpen(true);
}

export function registerModalApi(key: string, apiGetter: () => dialog.Api) {
	provider.set(key, apiGetter);
}
