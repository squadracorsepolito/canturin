<script lang="ts">
	import './styles.css';
	import * as editable from '@zag-js/editable';
	import { useMachine, normalizeProps, mergeProps } from '@zag-js/svelte';
	import type { EditableProps } from './types';

	let {
		value = $bindable(),
		name,
		placeholder,
		errors,
		textSize = 'md',
		fontWeight = 'normal',
		border = 'visible',
		readOnly = false,
		label,
		oncommit
	}: EditableProps<string> = $props();

	let fallbackValue = $state(value);

	const id = $props.id();
	const editableProps: editable.Props = $derived({
		id,
		name: name,
		value: value,
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

			const tmpValue = details.value;
			if (fallbackValue === tmpValue) {
				return;
			}
			fallbackValue = tmpValue;

			oncommit?.(details.value);
		},
		onValueChange: (details) => {
			value = details.value;
		}
	});

	const service = useMachine(editable.machine, () => editableProps);

	const api = $derived(editable.connect(service, normalizeProps));

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

<div class="editable">
	{#if label}
		<div>
			<label {...api.getLabelProps()}>
				<div>{label}</div>
			</label>
		</div>
	{/if}

	<div {...rootProps}>
		<div
			{...api.getAreaProps()}
			data-error={errors ? true : undefined}
			data-text-size={textSize}
			data-font-weight={fontWeight}
			data-border={border}
			data-readonly={readOnly ? true : undefined}
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
