// This file dynamically imports the correct environment file based on the NODE_ENV variable.

import { environment as localEnvironment } from './environment.local';
import { environment as developEnvironment } from './environment.develop';

const env = process.env.NODE_ENV || 'local';

export const environment = env === 'develop' ? developEnvironment : localEnvironment;
