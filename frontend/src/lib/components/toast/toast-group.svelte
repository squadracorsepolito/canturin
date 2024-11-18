<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import { normalizeProps, useMachine } from '@zag-js/svelte';
	import * as toast from '@zag-js/toast';
	import Toast from './toast.svelte';
	import { createToastProvider, pushToast } from './toast-provider.svelte';

	const [snapshot, send] = useMachine(
		toast.group.machine({
			id: uniqueId(),
			placement: 'bottom-end',
			duration: 3000,
			overlap: true,
			max: 5
		})
	);

	const api = $derived(toast.group.connect(snapshot, send, normalizeProps));

	createToastProvider(() => api);
</script>

{#each api.getPlacements() as placement}
	<div {...api.getGroupProps({ placement })}>
		{#each api.getToastsByPlacement(placement) as toast (toast.id)}
			<Toast actor={toast} />
		{/each}
	</div>
{/each}
