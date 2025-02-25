<script lang="ts">
	import './styles.css';
	import { uniqueId } from '$lib/utils';
	import * as editable from '@zag-js/editable';
	import * as numberInput from '@zag-js/number-input';
	import { useMachine, normalizeProps, mergeProps } from '@zag-js/svelte';
	import type { EditableProps } from './types';

	type Props = {
		min?: number;
		max?: number;
	};

	const inputId = uniqueId() + ':input';

	let {
		value = $bindable(),
		name,
		placeholder,
		errors,
		textSize = 'md',
		fontWeight = 'normal',
		border = 'visible',
		readOnly = false,
		oncommit,
		min,
		max
	}: EditableProps<number> & Props = $props();

	$inspect(min);

	let fallbackValue = $state(value);

	const editableProps: editable.Props = $derived({
		id: uniqueId(),
		name: name,
		value: value + '',
		readOnly: readOnly,
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
			if (fallbackValue === numValue) {
				return;
			}
			fallbackValue = numValue;

			oncommit?.(numValue);
		},
		onValueRevert(details) {
			value = +details.value;
		},
		ids: {
			input: inputId
		}
	});

	const editableService = useMachine(editable.machine, () => editableProps);

	const numberInputProps: numberInput.Props = $derived({
		id: uniqueId(),
		name: name,
		value: value + '',
		min: min,
		max: max,
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
	});

	const numberInputService = useMachine(numberInput.machine, () => numberInputProps);

	const api = $derived(editable.connect(editableService, normalizeProps));
	const inputApi = $derived(numberInput.connect(numberInputService, normalizeProps));

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
			data-readonly={readOnly ? true : undefined}
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
