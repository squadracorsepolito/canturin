export type EditableProps<T> = {
	value: T;
	name?: string;
	placeholder?: string;
	errors?: string[];
	textSize?: 'md' | 'lg';
	fontWeight?: 'normal' | 'medium' | 'semibold';
	border?: 'transparent' | 'visible';
	oncommit?: (value: T) => void;
};
