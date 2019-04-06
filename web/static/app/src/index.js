import { Application } from 'stimulus'
import { definitionsFromContext } from 'stimulus/webpack-helpers'
import './css/style.scss'
import { library, dom } from '@fortawesome/fontawesome-svg-core'
import { faCopy } from '@fortawesome/free-solid-svg-icons/faCopy'

library.add(faCopy)
dom.watch()
const application = Application.start()
const context = require.context('./controllers', true, /\.js$/)
application.load(definitionsFromContext(context))
