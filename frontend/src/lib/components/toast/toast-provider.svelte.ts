import * as toast from '@zag-js/toast';

class ToastProvider {
	#apiGetter: () => toast.GroupApi;

	constructor(apiGetter: () => toast.GroupApi) {
		this.#apiGetter = apiGetter;
	}

	push(type: toast.Type, title: string, description: string) {
		return this.#apiGetter().create({
			type,
			title,
			description
		});
	}
}

let provider: ToastProvider;

export function createToastProvider(apiGetter: () => toast.GroupApi) {
	provider = new ToastProvider(apiGetter);
}

export function pushToast(type: toast.Type, title: string, description: string) {
	if (!provider) {
		throw new Error('ToastProvider not initialized');
	}

	return provider.push(type, title, description);
}
