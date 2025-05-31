import moment from 'moment';
import moment_timezone from 'moment-timezone';

export default defineNuxtPlugin((nuxtApp) => {
  moment.locale("id");
  moment_timezone.tz.setDefault("Asia/Jakarta");

  nuxtApp.$moment = moment;
})