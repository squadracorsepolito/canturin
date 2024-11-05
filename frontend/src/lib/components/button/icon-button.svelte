<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		onclick: () => void;
		color?: 'primary' | 'error';
		label?: string;
		disabled?: boolean;
		children: Snippet;
	};

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
		@apply btn btn-sm btn-ghost;

		&:disabled {
			@apply opacity-75;
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
