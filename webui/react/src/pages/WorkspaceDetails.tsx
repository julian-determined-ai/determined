import { Tabs } from 'antd';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useHistory, useParams } from 'react-router';

import Page from 'components/Page';
import PageNotFound from 'components/PageNotFound';
import { useStore } from 'contexts/Store';
import { useFetchUsers } from 'hooks/useFetch';
import usePolling from 'hooks/usePolling';
import { paths } from 'routes/utils';
import { getWorkspace } from 'services/api';
import Message, { MessageType } from 'shared/components/Message';
import Spinner from 'shared/components/Spinner';
import { isNotFound } from 'shared/utils/service';
import { Workspace } from 'types';

import css from './WorkspaceDetails.module.scss';
import WorkspaceDetailsHeader from './WorkspaceDetails/WorkspaceDetailsHeader';
import WorkspaceMembers from './WorkspaceDetails/WorkspaceMembers';
import WorkspaceProjects from './WorkspaceDetails/WorkspaceProjects';

interface Params {
  tab: string;
  workspaceId: string;
}

export enum WorkspaceDetailsTab {
  Members = 'members',
  Projects = 'projects'
}

// Temporary Mock for rbacEnabled functionality
const rbacEnabled = false;

const WorkspaceDetails: React.FC = () => {
  const { users } = useStore();
  const { workspaceId } = useParams<Params>();
  const [ workspace, setWorkspace ] = useState<Workspace>();
  const [ pageError, setPageError ] = useState<Error>();
  const [ canceler ] = useState(new AbortController());
  const [ tabKey, setTabKey ] = useState<WorkspaceDetailsTab>(WorkspaceDetailsTab.Projects);
  const pageRef = useRef<HTMLElement>(null);
  const id = parseInt(workspaceId);
  const history = useHistory();
  const basePath = paths.workspaceDetails(workspaceId);
  const fetchWorkspace = useCallback(async () => {
    try {
      const response = await getWorkspace({ id }, { signal: canceler.signal });
      setWorkspace(response);
    } catch (e) {
      if (!pageError) setPageError(e as Error);
    }
  }, [ canceler.signal, id, pageError ]);

  const fetchUsers = useFetchUsers(canceler);

  const fetchAll = useCallback(async () => {
    await Promise.allSettled([ fetchWorkspace(), fetchUsers() ]);
  }, [ fetchWorkspace, fetchUsers ]);

  usePolling(fetchAll, { rerunOnNewFn: true });

  const handleTabChange = useCallback((activeTab) => {
    const tab = activeTab as WorkspaceDetailsTab;
    history.replace(`${basePath}/${tab}`);
    setTabKey(tab);
  }, [ basePath, history ]);

  useEffect(() => {
    return () => canceler.abort();
  }, [ canceler ]);

  if (isNaN(id)) {
    return <Message title={`Invalid Workspace ID ${workspaceId}`} />;
  } else if (pageError) {
    if (isNotFound(pageError)) return <PageNotFound />;
    const message =
      `Unable to fetch Workspace ${workspaceId}`;
    return <Message title={message} type={MessageType.Warning} />;
  } else if (!workspace) {
    return <Spinner tip={`Loading workspace ${workspaceId} details...`} />;
  }

  return (
    <Page
      className={css.base}
      containerRef={pageRef}
      headerComponent={(
        <WorkspaceDetailsHeader
          fetchWorkspace={fetchAll}
          workspace={workspace}
        />
      )}
      id="workspaceDetails">
      {
        rbacEnabled ? (
          <Tabs
            activeKey={tabKey}
            destroyInactiveTabPane
            onChange={handleTabChange}>
            <Tabs.TabPane
              destroyInactiveTabPane
              key={WorkspaceDetailsTab.Projects}
              tab="Projects">
              <WorkspaceProjects id={id} pageRef={pageRef} workspace={workspace} />
            </Tabs.TabPane>
            <Tabs.TabPane
              destroyInactiveTabPane
              key={WorkspaceDetailsTab.Members}
              tab="Members">
              <WorkspaceMembers pageRef={pageRef} users={users} workspace={workspace} />
            </Tabs.TabPane>
          </Tabs>
        ) : (<WorkspaceProjects id={id} pageRef={pageRef} workspace={workspace} />)
      }
    </Page>
  );
};

export default WorkspaceDetails;
