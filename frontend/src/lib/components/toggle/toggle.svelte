<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		toggled: boolean;
		name?: string;
		color?: 'primary' | 'error';
		children: Snippet;
	};

	let { toggled = $bindable(), name, color = 'primary', children }: Props = $props();
</script>

<button
	onclick={() => (toggled = !toggled)}
	aria-label={name}
	data-color={color}
	data-state={toggled ? 'on' : 'off'}
>
	{@render children()}
</button>

<style lang="postcss">
	button {
		@apply border-2 rounded-btn transition-colors p-1 shadow-none outline-none;

		&[data-state='off'] {
			@apply border-transparent;
		}

		&[data-color='primary'] {
			&[data-state='on'] {
				@apply border-primary bg-primary-ghost text-primary;
			}

			&[data-state='off'] {
				&:hover {
					@apply bg-primary-ghost text-primary;
				}
			}

			&:focus {
				@apply focus-ring-primary;
			}
		}

		&[data-color='error'] {
			&[data-state='on'] {
				@apply border-error bg-error-ghost text-error;
			}

			&[data-state='off'] {
				&:hover {
					@apply border-error bg-error-ghost text-error;
				}
			}

			&:focus {
				@apply focus-ring-error;
			}
		}
	}
</style>
