import * as toast from '@zag-js/toast';

class ToastProvider {
	#toaster: toast.Store;

	constructor(store: toast.Store) {
		this.#toaster = store;
	}

	push(type: toast.Type, title: string, description: string) {
		return this.#toaster.create({
			type,
			title,
			description
		});
	}
}

let provider: ToastProvider;

export function createToastProvider(store: toast.Store) {
	provider = new ToastProvider(store);
}

export function pushToast(type: toast.Type, title: string, description: string) {
	if (!provider) {
		throw new Error('ToastProvider not initialized');
	}

	return provider.push(type, title, description);
}
