<script lang="ts">
	import './styles.css';
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import * as numberInput from '@zag-js/number-input';
	import { useMachine, normalizeProps, mergeProps } from '@zag-js/svelte';
	import type { EditableProps } from './types';

	const inputId = uniqueId() + ':input';

	let {
		value = $bindable(),
		name,
		placeholder,
		errors,
		textSize = 'md',
		fontWeight = 'normal',
		border = 'visible',
		oncommit
	}: EditableProps<number> = $props();

	let fallbackValue = $state(value);

	const [snpshot, send] = useMachine(
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
			submitMode: 'both',
			autoResize: true,
			onValueCommit: (details) => {
				if (errors) {
					api.setValue(fallbackValue + '');
					value = fallbackValue;
					return;
				}

				const numValue = +details.value;

				fallbackValue = numValue;
				oncommit?.(numValue);
			},
			onValueRevert(details) {
				value = +details.value;
			},
			ids: {
				input: inputId
			}
		}),
		{
			context: {
				get value() {
					return value + '';
				}
			}
		}
	);

	const [inputSnapshot, inputSend] = useMachine(
		numberInput.machine({
			id: uniqueId(),
			name: name,
			onValueChange(details) {
				if (isNaN(details.valueAsNumber)) {
					value = 0;
					return;
				}

				value = details.valueAsNumber;
			},
			ids: {
				input: inputId
			}
		}),
		{
			context: {
				get value() {
					return value + '';
				}
			}
		}
	);

	const api = $derived(editable.connect(snpshot, send, normalizeProps));
	const inputApi = $derived(numberInput.connect(inputSnapshot, inputSend, normalizeProps));

	const rootProps = $derived(
		mergeProps(api.getRootProps(), {
			onkeydown: (e: KeyboardEvent) => {
				if (e.key === 'Escape') {
					api.cancel();
				}
			}
		})
	);
	const inputProps = $derived(mergeProps(inputApi.getInputProps(), api.getInputProps()));
</script>

<div class="editable">
	<div {...rootProps}>
		<div
			{...api.getAreaProps()}
			data-error={errors ? true : undefined}
			data-text-size={textSize}
			data-font-weight={fontWeight}
			data-border={border}
		>
			<input {...inputProps} />

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
