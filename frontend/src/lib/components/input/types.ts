export type InputProps<T> = {
	value: T;
	name: string;
	type: 'number' | 'text';
	size?: 'md' | 'sm';
	label?: string;
	width?: number;
	errors?: string[];
	focusOnDisplay?: boolean;
	min?: number;
	max?: number;
	onblur?: () => void;
	onescape?: () => void;
};

export type TextInputProps = Omit<InputProps<string>, 'type' | 'width' | 'min' | 'max'>;
export type ResizeableTextInputProps = TextInputProps;
export type NumberInputProps = Omit<InputProps<number>, 'type' | 'width'>;

export type RadioInputOption = {
	id: number;
	name: string;
	label: string;
	desc?: string;
	disabled?: boolean;
};
