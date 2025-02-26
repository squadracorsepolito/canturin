<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as dialog from '@zag-js/dialog';
	import { portal, normalizeProps, useMachine } from '@zag-js/svelte';
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		title: string;
		desc?: string;
		trigger: Snippet<[{ getProps: () => HTMLButtonAttributes }]>;
		content?: Snippet;
		actions?: Snippet<[{ close: () => void }]>;
	};

	let { title, desc, trigger, content, actions }: Props = $props();

	const service = useMachine(dialog.machine, {
		id: uniqueId()
	});

	const api = $derived(dialog.connect(service, normalizeProps));
</script>

{@render trigger({ getProps: api.getTriggerProps })}

{#if api.open}
	<div use:portal {...api.getBackdropProps()}></div>

	<div use:portal {...api.getPositionerProps()}>
		<div {...api.getContentProps()}>
			<h3 {...api.getTitleProps()}>{title}</h3>

			{#if desc}
				<p {...api.getDescriptionProps()}>{desc}</p>
			{/if}

			{#if content}
				<div class="pt-6">
					{@render content()}
				</div>
			{/if}

			<div class="modal-action">
				{@render actions?.({ close: () => api.setOpen(false) })}

				<button {...api.getCloseTriggerProps()}>Cancel</button>
			</div>
		</div>
	</div>
{/if}

<style lang="postcss">
	[data-scope='dialog'][data-part='backdrop'] {
		@apply bg-neutral opacity-75 block fixed w-full h-full top-0 left-0;
	}

	[data-scope='dialog'][data-part='positioner'] {
		@apply fixed w-full h-full top-0 left-0 flex items-center justify-center;
	}

	[data-scope='dialog'][data-part='content'] {
		@apply modal-box;

		[data-part='description'] {
			@apply opacity-85 pt-1;
		}

		[data-part='close-trigger'] {
			@apply btn;
		}
	}
</style>
