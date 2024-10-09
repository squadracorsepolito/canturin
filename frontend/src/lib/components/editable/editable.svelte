<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { z } from 'zod';

	type Props = {
		initialValue: string;
		validator: z.AnyZodObject;
		placeholder?: string;
		onSubmit: (value: string) => void;
	};

	let { initialValue, validator, placeholder, onSubmit }: Props = $props();

	let errors = $state<string[]>();

	const [snpshot, send] = useMachine(
		editable.machine({
			id: uniqueId(),
			value: initialValue,
			activationMode: 'dblclick',
			placeholder: placeholder,
			autoResize: true,
			submitMode: 'both',
			onValueCommit: (details) => {
				if (errors) {
					api.setValue(initialValue);
					errors = undefined;
					return;
				}

				onSubmit(details.value);
			},
			onValueChange: (details) => {
				const res = validator.safeParse({ name: details.value });
				if (!res.success) {
					errors = res.error.flatten().fieldErrors['name'];
				} else {
					errors = undefined;
				}
			}
		})
	);

	const api = $derived(editable.connect(snpshot, send, normalizeProps));
</script>

<div {...api.getRootProps()}>
	<div {...api.getAreaProps()} data-error={errors ? true : undefined}>
		<input {...api.getInputProps()} />

		<span {...api.getPreviewProps()}>
			{api.valueText}
		</span>
	</div>
</div>

{#if errors}
	<div class="absolute pt-1 text-error text-xs">
		{#each errors as err}
			<span>{err}</span>
		{/each}
	</div>
{/if}

<style lang="postcss">
	[data-part='area'] {
		@apply rounded-btn border-2 border-transparent px-3 py-1 font-medium text-xl;

		&[data-error] {
			&[data-focus] {
				@apply focus-ring-error border-error;
			}
		}

		&:not([data-error]) {
			&[data-focus] {
				@apply focus-ring-primary border-primary;
			}
		}
	}
</style>
