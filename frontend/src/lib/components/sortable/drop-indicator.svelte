<script lang="ts">
	import type { Edge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/dist/types/types';

	type Props = {
		edge: Edge;
	};

	let { edge }: Props = $props();
</script>

<div class="drop-indicator" data-edge={edge}></div>

<style lang="postcss">
	.drop-indicator {
		@apply absolute h-1 w-full bg-accent z-10 rounded-r-badge;

		&::before {
			content: '';

			@apply absolute w-1 h-4 bg-accent rounded-badge -left-1 top-[-0.375rem];
		}

		&[data-edge='top'] {
			@apply top-[-0.25rem];
		}

		&[data-edge='bottom'] {
			@apply bottom-[-0.25rem];
		}
	}
</style>
