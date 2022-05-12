import 'styles/index.scss';
import 'styles/storybook.scss';
import 'shared/prototypes';

import ThemeDecorator from "storybook/ThemeDecorator"

export const decorators = [ ThemeDecorator ];
export const parameters = { layout: 'centered' };
