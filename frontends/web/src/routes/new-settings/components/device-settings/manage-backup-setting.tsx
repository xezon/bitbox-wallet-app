/**
 * Copyright 2023 Shift Crypto AG
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useTranslation } from 'react-i18next';
import { SettingsItem } from '../settingsItem/settingsItem';
import { ChevronRightDark } from '../../../../components/icon';
import { route } from '../../../../utils/route';

type TProps = {
  deviceID: string;
}

const ManageBackupSetting = ({ deviceID }: TProps) => {
  const { t } = useTranslation();
  return (
    <SettingsItem
      onClick={() => route(`/manage-backups/${deviceID}`)}
      settingName={t('backup.title')}
      secondaryText={t('deviceSettings.backups.manageBackups.description')}
      extraComponent={<ChevronRightDark />}
    />
  );
};


export { ManageBackupSetting };