import { useCallback, useState } from 'react';

import usePolling from 'hooks/usePolling';
import { queryTrials } from 'services/api';
import handleError from 'utils/error';

import { encodeFilters, encodeTrialSorter } from '../api';
import { TrialFilters, TrialSorter } from '../Collections/filters';

import { decodeTrialsWithMetadata, defaultTrialData, TrialsWithMetadata } from './data';
import {V1Pagination} from 'services/api-ts-sdk';
import { number } from 'fp-ts';

interface Params {
  filters: TrialFilters;
  limit: number;
  offset: number;
  sorter: TrialSorter;
}

interface TrialsWithMetadataWithPagination {
  trials: TrialsWithMetadata;
  pagination: V1Pagination;
  total: number;
}

export const useFetchTrials = ({
  filters,
  limit,
  offset,
  sorter,
}: Params): TrialsWithMetadataWithPagination => {
  const [ trials, setTrials ] = useState<TrialsWithMetadata>(defaultTrialData());
  const [ pagination, setPagination ] = useState<V1Pagination>({});
  const [ total, setTotal ] = useState<number>(0);
  const fetchTrials = useCallback(async () => {
    let response: any;
    const _filters = encodeFilters(filters);
    const _sorter = encodeTrialSorter(sorter);
    try {
      response = await queryTrials({
        filters: _filters,
        pagination: {
          limit,
          offset
        },
        sorter: _sorter,
      });
    } catch (e) {
      handleError(e, { publicSubject: 'Unable to fetch trials.' });
    }
    if (response){
      const newTrials = decodeTrialsWithMetadata(response.trials);
      setPagination(response.pagination);
      setTotal(response.total);
      if (newTrials)
        setTrials(newTrials);
    }

  }, [ filters, limit, offset, sorter ]);

  usePolling(fetchTrials, { interval: 200000, rerunOnNewFn: true });
  return {trials, pagination, total} ;
};
