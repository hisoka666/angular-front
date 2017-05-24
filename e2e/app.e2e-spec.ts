import { AngularFrontPage } from './app.po';

describe('angular-front App', () => {
  let page: AngularFrontPage;

  beforeEach(() => {
    page = new AngularFrontPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
