import { Application } from 'stimulus'
import { definitionsFromContext } from 'stimulus/webpack-helpers'
import '../node_modules/toastr/build/toastr.css'
import './css/style.scss'
import { library, dom } from '@fortawesome/fontawesome-svg-core'
import { faCopy } from '@fortawesome/free-solid-svg-icons/faCopy'
import { faCog } from '@fortawesome/free-solid-svg-icons/faCog'

library.add(faCog)
dom.watch()

library.add(faCopy)
dom.watch()
const application = Application.start()
const context = require.context('./controllers', true, /\.js$/)
application.load(definitionsFromContext(context))
