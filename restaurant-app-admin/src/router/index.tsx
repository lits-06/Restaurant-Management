// import { createBrowserRouter, RouterProvider } from 'react-router-dom'
// import MainLayout from '../components/layout/MainLayout'
// import HomePage from '../pages/HomePage'

// /**
//  * Các page khác sẽ được import tại đây khi build:
//  *   import MenuPage    from '../pages/MenuPage'
//  *   import BookingPage from '../pages/BookingPage'
//  *   import CheckoutPage from '../pages/CheckoutPage'
//  *   import LoginPage   from '../pages/LoginPage'
//  */

// const router = createBrowserRouter([
//   {
//     path: '/',
//     element: <MainLayout />,
//     children: [
//       { index: true,     element: <HomePage /> },
//       // { path: 'menu',    element: <MenuPage /> },
//       // { path: 'booking', element: <BookingPage /> },
//       // { path: 'checkout',element: <CheckoutPage /> },
//       // { path: 'login',   element: <LoginPage /> },
//       // { path: 'about',   element: <AboutPage /> },
//     ],
//   },
// ])

// export default function AppRouter() {
//   return <RouterProvider router={router} />
// }