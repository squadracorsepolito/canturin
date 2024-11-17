<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import { normalizeProps, useMachine } from '@zag-js/svelte';
	import * as toast from '@zag-js/toast';
	import Toast from './toast.svelte';

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
</script>

<!-- <button
	onclick={() =>
		api.create({
			type: 'error',
			title: 'Error',
			description:
				'Something went wrong Something went wrong Something went wrong Something went wrong Something went wrong Something went wrong'
		})}
>
	add
</button> -->

{#each api.getPlacements() as placement}
	<div {...api.getGroupProps({ placement })}>
		{#each api.getToastsByPlacement(placement) as toast (toast.id)}
			<Toast actor={toast} />
		{/each}
	</div>
{/each}
