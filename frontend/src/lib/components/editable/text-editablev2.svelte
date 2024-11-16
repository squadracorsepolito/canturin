<script lang="ts">
	import './styles.css';
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import { useMachine, normalizeProps, mergeProps } from '@zag-js/svelte';

	type Props = {
		value: string;
		name?: string;
		placeholder?: string;
		errors?: string[];
		textSize?: 'md' | 'lg';
		fontWeight?: 'normal' | 'medium';
		oncommit?: (value: string) => void;
	};

	let {
		value = $bindable(),
		name,
		placeholder,
		errors,
		textSize = 'md',
		fontWeight = 'normal',
		oncommit
	}: Props = $props();

	let fallbackValue = $state(value);

	const [snapshot, send] = useMachine(
		editable.machine({
			id: uniqueId(),
			name: name,
			activationMode: 'dblclick',
			placeholder: placeholder
				? {
						edit: '',
						preview: placeholder
					}
				: undefined,
			autoResize: true,
			submitMode: 'both',
			onValueCommit: (details) => {
				if (errors) {
					api.setValue(fallbackValue);
					return;
				}

				fallbackValue = details.value;
				oncommit?.(details.value);
			},
			onValueChange: (details) => {
				value = details.value;
			}
		}),
		{
			context: {
				get value() {
					return value;
				}
			}
		}
	);

	const api = $derived(editable.connect(snapshot, send, normalizeProps));

	const rootProps = $derived(
		mergeProps(api.getRootProps(), {
			onkeydown: (e: KeyboardEvent) => {
				if (e.key === 'Escape') {
					api.cancel();
				}
			}
		})
	);
</script>

<div class="relative editable">
	<div {...rootProps}>
		<div
			{...api.getAreaProps()}
			data-error={errors ? true : undefined}
			data-text-size={textSize}
			data-font-weight={fontWeight}
		>
			<input {...api.getInputProps()} />

			<span {...api.getPreviewProps()}>
				{api.valueText}
			</span>
		</div>
	</div>

	{#if errors && api.editing}
		<div data-part="error">
			{#each errors as err}
				<span>{err}</span>
			{/each}
		</div>
	{/if}
</div>
