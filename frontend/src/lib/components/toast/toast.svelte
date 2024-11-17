<script lang="ts">
	import { normalizeProps, useActor } from '@zag-js/svelte';
	import * as toast from '@zag-js/toast';
	import { CloseIcon } from '../icon';

	type Props = {
		actor: toast.Service;
	};

	const { actor }: Props = $props();

	const [snapshot, send] = useActor(actor);

	const api = $derived(toast.connect(snapshot, send, normalizeProps));

	$inspect(actor.config);
</script>

<div {...api.getRootProps()}>
	<div {...api.getGhostBeforeProps()}></div>

	<div>
		<h4 {...api.getTitleProps()}>{api.title}</h4>

		<p {...api.getDescriptionProps()}>{api.description}</p>
	</div>

	<button {...api.getCloseTriggerProps()} class="btn btn-sm btn-ghost btn-square">
		<CloseIcon />
	</button>

	<div {...api.getGhostAfterProps()}></div>
</div>

<style lang="postcss">
	[data-part='root'] {
		translate: var(--x) var(--y);
		scale: var(--scale);
		z-index: var(--z-index);
		height: var(--height);
		opacity: var(--opacity);
		will-change: translate, opacity, scale;

		@apply min-w-96 p-5 gap-3  rounded-box flex transition-all;

		&[data-type='error'] {
			@apply bg-error text-error-content;
		}

		[data-part='title'] {
			@apply font-semibold;
		}

		[data-part='description'] {
			@apply text-sm pt-1 line-clamp-3;
		}
	}
</style>
