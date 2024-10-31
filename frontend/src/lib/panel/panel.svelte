<script lang="ts">
	import type { Snippet } from 'svelte';
	import { autoScrollForElements } from '@atlaskit/pragmatic-drag-and-drop-auto-scroll/element';
	import type { Action } from 'svelte/action';

	let { children }: { children: Snippet } = $props();

	const panelAction: Action<HTMLElement> = (el) => {
		const cleanup = autoScrollForElements({
			element: el
		});

		return {
			destroy() {
				cleanup();
			}
		};
	};
</script>

<article use:panelAction class="flex-1 overflow-y-auto p-5 max-w-none">
	<div class="container m-auto @container flex flex-col gap-8">
		{@render children()}
	</div>
</article>
