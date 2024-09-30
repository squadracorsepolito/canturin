export type InputProps<T> = {
	value: T;
	name: string;
	type: 'number' | 'text';
	label?: string;
	errors?: string[];
	focusOnDiplay?: boolean;
	onblur?: () => void;
	onescape?: () => void;
};
