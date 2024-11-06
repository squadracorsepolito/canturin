<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = {
		onclick: () => void;
		color?: 'primary' | 'error';
		label?: string;
		disabled?: boolean;
		children: Snippet;
	} & HTMLButtonAttributes;

	let { onclick, color, label, disabled, children }: Props = $props();
</script>

<button {onclick} class="icon-button {!label && 'btn-square'}" {disabled} data-color={color}>
	{@render children()}

	{#if label}
		{label}
	{/if}
</button>

<style lang="postcss">
	.icon-button {
		@apply btn btn-sm btn-ghost border-2 py-1 h-auto;

		&:disabled {
			@apply opacity-75 cursor-not-allowed;
		}

		&[data-color='primary'] {
			@apply text-primary;

			&:hover,
			&:disabled {
				@apply bg-primary-ghost;
			}
		}

		&[data-color='error'] {
			@apply text-error;

			&:hover,
			&:disabled {
				@apply bg-error-ghost;
			}
		}
	}
</style>
