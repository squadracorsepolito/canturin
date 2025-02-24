<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import { normalizeProps, useMachine } from '@zag-js/svelte';
	import * as toast from '@zag-js/toast';
	import Toast from './toast.svelte';
	import { createToastProvider } from './toast-provider.svelte';

	const store = toast.createStore({
		placement: 'bottom-end',
		duration: 3000,
		max: 5,
		overlap: true
	});

	const service = useMachine(toast.group.machine, {
		id: uniqueId(),
		store: store
	});

	const api = $derived(toast.group.connect(service, normalizeProps));

	createToastProvider(store);
</script>

<div {...api.getGroupProps()}>
	{#each api.getToasts() as toast, index (toast.id)}
		<Toast {toast} {index} parent={service} />
	{/each}
</div>
