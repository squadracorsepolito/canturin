<script lang="ts">
	import type { ThemeColorProps } from '$lib/utils/types';
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = HTMLButtonAttributes &
		ThemeColorProps & {
			label?: string;
			children?: Snippet;
		};

	let { onclick, themeColor, label, disabled, children, ...rest }: Props = $props();
</script>

<button
	{onclick}
	{...rest}
	class="icon-button {!label && 'btn-square'}"
	{disabled}
	data-theme-color={themeColor}
>
	{@render children?.()}

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

		&[data-theme-color='primary'] {
			@apply text-primary;

			&:hover,
			&:disabled {
				@apply bg-primary-ghost;
			}
		}

		&[data-theme-color='secondary'] {
			@apply text-secondary;

			&:hover,
			&:disabled {
				@apply bg-secondary-ghost;
			}
		}

		&[data-theme-color='error'] {
			@apply text-error;

			&:hover,
			&:disabled {
				@apply bg-error-ghost;
			}
		}
	}
</style>
