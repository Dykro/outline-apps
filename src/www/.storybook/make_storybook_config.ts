/*
  Copyright 2022 The Outline Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

import {PolymerElement} from "@polymer/polymer";
import {LegacyElementMixin} from "@polymer/polymer/lib/legacy/legacy-element-mixin";

class LegacyPolymerElement extends LegacyElementMixin(PolymerElement) {}

interface StorybookGenericControl<T> {
  defaultValue?: T;
  controlName: string;
}

interface StorybookTextControl extends StorybookGenericControl<string> {
  controlType: "text" | "color" | "date";
}

interface StorybookObjectControl extends StorybookGenericControl<object> {
  controlType: "object";
}

interface StorybookSelectControl extends StorybookGenericControl<string | number | string[] | number[]> {
  controlType: "radio" | "inline-radio" | "check" | "inline-check" | "select" | "multi-select";
  options: string[] | number[];
}

interface StorybookBooleanControl extends StorybookGenericControl<boolean> {
  controlType: "boolean";
}

export type StorybookControl =
  | StorybookTextControl
  | StorybookSelectControl
  | StorybookObjectControl
  | StorybookBooleanControl;

interface MakeStorybookConfigOptions {
  containerPath?: string;
  controls: StorybookControl[];
}

interface StorybookConfig {
  name: string;
  component: string;
  args: {[argName: string]: string | object | boolean | number | string[] | number[]};
  argTypes: {[argName: string]: {control: string; options?: string[] | number[]}};
}

export function makeStorybookConfig(
  Component: LegacyPolymerElement,
  {controls, containerPath: containerName}: MakeStorybookConfigOptions
): StorybookConfig {
  const componentName = Component.constructor.name;

  const result: StorybookConfig = {
    name: containerName ? `${containerName}/${componentName}` : componentName,
    component: Component.is,
    args: {},
    argTypes: {},
  };

  for (const control of controls) {
    result.args[control.controlName] = control.defaultValue;
    result.argTypes[control.controlName] = {control: control.controlType};

    if ("options" in control) {
      result.argTypes[control.controlName].options = control.options;
    }
  }

  return result;
}
