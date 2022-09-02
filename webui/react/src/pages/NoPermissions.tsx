import React, { useEffect } from 'react';

import Page from 'components/Page';
import { StoreAction, useStore, useStoreDispatch } from 'contexts/Store';
import Icon from 'shared/components/Icon/Icon';

import css from './NoPermissions.module.scss';

const NoPermissions: React.FC = () => {

  const { users } = useStore();

  const systemAdmin = users.find((user) => user.isAdmin);
  // TODO: Implement querying for the system admin

  const storeDispatch = useStoreDispatch();
  useEffect(() => storeDispatch({ type: StoreAction.HideUIChrome }), [ storeDispatch ]);

  return (
    <Page
      bodyNoPadding>
      <div className={css.base}>
        <div className={css.icon}>
          <Icon name="warning-large" size="mega" />
        </div>
        <h6>You don&lsquo;t have access to a workspace</h6>
        <p className={css.message}>
          In order to access and use Determined you need to be added to a workspace.
          Contact {systemAdmin?.username } to request to be added to a workspace.
        </p>
      </div>
    </Page>
  );
};

export default NoPermissions;
