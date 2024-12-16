// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { useEffect, useState } from 'react';

interface SampleDataModule<T> {
  default: T;
}

export function useSampleData<T>(name: string): T | undefined {
  const [data, setData] = useState<T>();
  useEffect(() => {
    async function loadData(): Promise<void> {
      const js: SampleDataModule<T> = await import(`../../sample-data/${name}`);
      setData(js.default);
    }
    loadData();
  }, [name]);
  return data;
}
