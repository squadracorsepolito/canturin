<script lang="ts">
	import { normalizeProps, useMachine } from '@zag-js/svelte';
	import * as splitter from '@zag-js/splitter';
	import type { Snippet } from 'svelte';

	type PanelProps = {
		size?: number;
		minSize?: number;
		maxSize?: number;
	};

	type Props = {
		left: Snippet;
		leftPanel?: PanelProps;
		right: Snippet;
		rightPanel?: PanelProps;
	};

	let { left, leftPanel, right, rightPanel }: Props = $props();

	const id = $props.id();
	const service = useMachine(splitter.machine, {
		id,
		defaultSize: [
			{
				id: 'left',
				size: leftPanel?.size,
				minSize: leftPanel?.minSize,
				maxSize: leftPanel?.maxSize
			},
			{
				id: 'right',
				size: rightPanel?.size,
				minSize: rightPanel?.minSize,
				maxSize: rightPanel?.maxSize
			}
		]
	});

	const api = $derived(splitter.connect(service, normalizeProps));
</script>

<div {...api.getRootProps()}>
	<div {...api.getPanelProps({ id: 'left' })}>
		{@render left()}
	</div>

	<div {...api.getResizeTriggerProps({ id: 'left:right' })}></div>

	<div {...api.getPanelProps({ id: 'right' })}>
		{@render right()}
	</div>
</div>

<style lang="postcss">
	:root {
		[data-scope='splitter'][data-part='resize-trigger'] {
			@apply w-1 bg-base-300 transition-colors;

			&[data-focus] {
				@apply bg-accent;
			}

			&:active {
				@apply focus-ring-accent;
			}
		}
	}
</style>
